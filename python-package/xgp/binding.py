from ctypes import *

import numpy as np


class GoSlice(Structure):
    _fields_ = [
        ('data', POINTER(c_float)),
        ('len', c_longlong),
        ('cap', c_longlong)
    ]


class GoMatrix(Structure):
    _fields_ = [
        ('data', POINTER(GoSlice)),
        ('len', c_longlong),
        ('cap', c_longlong)
    ]


def numpy_to_go(arr: np.ndarray) -> GoSlice:
    # If the slice is 1D then return a GoSlice
    if len(arr.shape) == 1:
        return GoSlice(
            arr.ctypes.data_as(POINTER(c_float)),
            len(arr),
            len(arr)
        )
    return GoMatrix(
        (GoSlice * len(arr))(*[numpy_to_go(row) for row in arr]),
        len(arr),
        len(arr)
    )

def fit(X: np.ndarray, y: np.ndarray, metric_name: str, generations: int,
        tuning_generations: int):
    """Refers to the Fit method in main.go"""
    xgp = cdll.LoadLibrary('./xgp.so')
    xgp.Fit.argtypes = [
        GoMatrix,
        GoSlice,
        c_wchar_p,
        c_int,
        c_int
    ]
    print(generations, tuning_generations, c_int(generations), c_int(tuning_generations))
    return xgp.Fit(
        numpy_to_go(X),
        numpy_to_go(y),
        metric_name,
        c_int(generations),
        c_int(tuning_generations)
    )


def predict(X: np.ndarray) -> np.ndarray:
    """Refers to the Predict method in main.go"""
    xgp = cdll.LoadLibrary('./xgp.so')
    xgp.Predict.argtypes = [GoMatrix]
    xgp.Predict.restype = GoSlice
    return xgp.Predict(numpy_to_go(X))
