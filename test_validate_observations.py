import json, unittest
from pathlib import Path
from validate_observations import validate
ROOT=Path(__file__).resolve().parents[1]
class ValidationTests(unittest.TestCase):
    def setUp(self): self.data=json.loads((ROOT/'examples/observations.json').read_text())
    def test_valid(self): self.assertEqual(validate(self.data), [])
    def test_scope_violation(self):
        self.data['observations'][0]['authorized_scope_id']='other'
        self.assertTrue(any('out-of-scope' in e for e in validate(self.data)))
    def test_unsubstantiated_inference(self):
        self.data['observations'][1]['evidence_refs']=[]
        self.assertTrue(any('requires evidence_refs' in e for e in validate(self.data)))
if __name__=='__main__': unittest.main()
