import unittest

import numpy as np

from koza import binding


class TestBinding(unittest.TestCase):

    def test_numpy_to_float64_slice(self):
        """Test numpy_to_float64_slice"""
        X = np.array([
            [11, 12, 13],
            [21, 22, 23]
        ])
        go_slice = binding.numpy_to_float64_slice(X)
        assert isinstance(go_slice, binding.GoFloat64Matrix)
