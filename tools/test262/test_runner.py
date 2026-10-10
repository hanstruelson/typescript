import importlib.util
from pathlib import Path
import tempfile
import unittest

SPEC = importlib.util.spec_from_file_location('test262_runner', Path(__file__).with_name('run.py'))
runner = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(runner)

class InventoryTests(unittest.TestCase):
    def test_modes_exclusions_and_imported_fixture_files(self):
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder)
            (root/'test/language/expressions/addition').mkdir(parents=True)
            directory = root/'test/language/expressions/addition'
            sources = {
                'both.js': '/*---\ndescription: addition\n---*/\n1+2;',
                'strict.js': '/*---\nflags: [onlyStrict]\n---*/\n1;',
                'sloppy.js': '/*---\nflags: [noStrict]\n---*/\n1;',
                'module.js': '/*---\nflags: [module]\n---*/\nexport {};',
                'dynamic.js': '/*---\ndescription: dynamic\n---*/\neval("1");',
                'ordinary-function.js': '/*---\ndescription: prototype\n---*/\nFunction.prototype;',
                'fixture.js': 'export const value=1;',
                'negative.js': '/*---\nnegative:\n  phase: parse\n  type: SyntaxError\n---*/\nconst =;',
            }
            for name, source in sources.items(): (directory/name).write_text(source)
            cases = {Path(c['path']).name:c for c in runner.inventory(root)}
            self.assertNotIn('fixture.js', cases)
            self.assertEqual(cases['both.js']['modes'], ['sloppy','strict'])
            self.assertEqual(cases['strict.js']['modes'], ['strict'])
            self.assertEqual(cases['sloppy.js']['modes'], ['sloppy'])
            self.assertEqual(cases['module.js']['modes'], ['module'])
            self.assertEqual(cases['dynamic.js']['excluded'], 'explicit-dynamic-code')
            self.assertIsNone(cases['ordinary-function.js']['excluded'])
            self.assertEqual(cases['negative.js']['metadata']['negative']['phase'], 'parse')

if __name__=='__main__': unittest.main()

class OutcomeTests(unittest.TestCase):
    def run_case(self, responses, negative=None):
        from types import SimpleNamespace
        from unittest.mock import patch
        with tempfile.TemporaryDirectory() as folder:
            root=Path(folder)
            (root/'test').mkdir()
            (root/'test/example.js').write_text('/*---\nflags: [raw]\n---*/\n1;')
            args=SimpleNamespace(suite=root,node='node',go='go',emitter=Path('/fake/emitter'),timeout=1,build_timeout=1)
            meta={'flags':['raw']}
            if negative: meta['negative']=negative
            case={'path':'example.js','category':'language/example','metadata':meta}
            with patch.object(runner,'command',side_effect=responses):
                return runner.execute(case,'raw',args)['status']

    def test_native_build_error_is_not_pass(self):
        self.assertEqual(self.run_case([(0,runner.MARK,0),(0,'',0),(1,'invalid Go',0)]),'build-error')

    def test_expected_parse_error_requires_syntactic_rejection(self):
        negative={'phase':'parse','type':'SyntaxError'}
        self.assertEqual(self.run_case([(1,'SyntaxError',0),(3,'syntax rejected',0)],negative),'pass')
        self.assertEqual(self.run_case([(1,'SyntaxError',0),(4,'unsupported emission',0)],negative),'emit-error')

    def test_successful_exit_without_completion_is_not_pass(self):
        self.assertEqual(self.run_case([(0,runner.MARK,0),(0,'',0),(0,'',0),(0,'',0)]),'missing-completion')

    def test_reference_failure_does_not_skip_native_build(self):
        self.assertEqual(self.run_case([(1,'missing reference feature',0),(0,'',0),(1,'invalid Go',0)]),'build-error')

class PrivateCacheTests(unittest.TestCase):
    def test_prunes_new_go127_directories_and_preserves_seed(self):
        import os
        from unittest.mock import patch
        with tempfile.TemporaryDirectory() as folder:
            base=Path(folder);seed=base/'go-build-cache';cache=base/'test262-go-cache'
            seed_entry=seed/'ab/seed-d';seed_entry.mkdir(parents=True)
            (seed_entry/'archive.a').write_bytes(b'shared cache')
            cached_seed=cache/'ab/seed-d';cached_seed.mkdir(parents=True)
            os.link(seed_entry/'archive.a',cached_seed/'archive.a')
            fresh=cache/'ab/new-d';fresh.mkdir();(fresh/'archive.a').write_bytes(b'task cache')
            os.utime(fresh,(0,0));os.utime(cached_seed,(0,0))
            with patch.dict(os.environ,{'GOCACHE':str(cache),'TEST262_PRIVATE_GOCACHE':str(cache),'TEST262_CACHE_SEED':str(seed)}):
                runner.trim_private_cache()
            self.assertFalse(fresh.exists())
            self.assertEqual((cached_seed/'archive.a').read_bytes(),b'shared cache')
            self.assertEqual((seed_entry/'archive.a').read_bytes(),b'shared cache')
