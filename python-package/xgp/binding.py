from ctypes import *

import numpy as np


class GoString(Structure):
    _fields_ = [
        ('p', c_char_p),
        ('n', c_longlong)
    ]


class GoStringSlice(Structure):
    _fields_ = [
        ('data', POINTER(c_wchar_p)),
        ('len', c_longlong),
        ('cap', c_longlong)
    ]


class GoFloat64Slice(Structure):
    _fields_ = [
        ('data', POINTER(c_double)),
        ('len', c_longlong),
        ('cap', c_longlong)
    ]


class GoFloat64Matrix(Structure):
    _fields_ = [
        ('data', POINTER(GoFloat64Slice)),
        ('len', c_longlong),
        ('cap', c_longlong)
    ]


def numpy_to_float64_slice(arr: np.ndarray) -> GoFloat64Slice:
    # If the slice is 1D then return a GoFloat64Slice
    if len(arr.shape) == 1:
        return GoFloat64Slice(
            arr.ctypes.data_as(POINTER(c_double)),
            len(arr),
            len(arr)
        )
    return GoFloat64Matrix(
        (GoFloat64Slice * len(arr))(*[numpy_to_float64_slice(row) for row in arr]),
        len(arr),
        len(arr)
    )


def str_list_to_string_slice(l: list) -> GoStringSlice:
    return GoStringSlice((c_wchar_p * len(l))(*l), len(l), len(l))


def fit(X: np.ndarray,
        y: np.ndarray,
        X_names: list,
        const_max: float,
        const_min: float,
        eval_metric_name: str,
        funcs_string: str,
        generations: int,
        loss_metric_name: str,
        max_height: int,
        min_height: int,
        n_pops: int,
        parsimony_coeff: float,
        p_constant: float,
        p_terminal: float,
        pop_size: int,
        rounds: int,
        tuning_generations: int,
        seed: int,
        verbose: bool):
    """Refers to the Fit method in main.go"""
    xgp = cdll.LoadLibrary('./xgp.so')
    xgp.Fit.argtypes = [
        GoFloat64Matrix, # X
        GoFloat64Slice, # y
        GoStringSlice, # X_names
        c_double, # const_max
        c_double, # const_min
        GoString, # eval_metric_name
        GoString, # funcs_string
        c_longlong, # generations
        GoString, # loss_metric_name
        c_longlong, # max_height
        c_longlong, # min_height
        c_longlong, # n_pops
        c_double, # parsimony_coeff
        c_double, # p_constant
        c_double, # p_terminal
        c_longlong, # pop_size
        c_longlong, # rounds
        c_longlong, # tuning_generations
        c_longlong, # seed
        c_bool # verbose
    ]
    xgp.Fit.restype = c_char_p

    program_bytes = xgp.Fit(
        numpy_to_float64_slice(np.transpose(X)),
        numpy_to_float64_slice(y),
        str_list_to_string_slice(X_names),
        const_max,
        const_min,
        GoString(bytes(eval_metric_name, 'utf-8'), len(eval_metric_name)),
        GoString(bytes(funcs_string, 'utf-8'), len(funcs_string)),
        generations,
        GoString(bytes(loss_metric_name, 'utf-8'), len(loss_metric_name)),
        max_height,
        min_height,
        n_pops,
        parsimony_coeff,
        p_constant,
        p_terminal,
        pop_size,
        rounds,
        tuning_generations,
        seed,
        verbose
    )

    return program_bytes.decode()


# def predict(X: np.ndarray, predict_proba: bool) -> np.ndarray:
#     """Refers to the Predict method in main.go"""
#     xgp = cdll.LoadLibrary('./xgp.so')
#     xgp.Predict.argtypes = [GoFloat64Matrix, c_bool]
#     #xgp.Predict.restype = GoFloat64Slice
#     y_pred = xgp.Predict(numpy_to_float64_slice(np.transpose(X)), predict_proba)
#     print(y_pred)
#     return y_pred
