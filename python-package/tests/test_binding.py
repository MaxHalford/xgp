import unittest

import numpy as np

from xgp import binding


class TestBinding(unittest.TestCase):

    def test_np_to_go(self):
        """Test array_to_go_slice works with np.ndarrays."""
        X = np.array([
            [11, 12, 13],
            [21, 22, 23]
        ])
        go_slice = binding.array_to_go_slice(X)
        assert isinstance(go_slice, GoSlice)
