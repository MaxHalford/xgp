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
        n_generations: int,
        loss_metric_name: str,
        max_height: int,
        min_height: int,
        n_populations: int,
        p_constant: float,
        p_full: float,
        p_hoist_mutation: float,
        p_point_mutation: float,
        p_subtree_crossover: float,
        p_subtree_mutation: float,
        p_terminal: float,
        parsimony_coeff: float,
        point_mutation_rate: float,
        population_size: int,
        n_rounds: int,
        tuning_n_generations: int,
        seed: int,
        verbose: bool):
    """Refers to the Fit method in main.go"""
    koza = cdll.LoadLibrary('./koza.so')
    koza.Fit.argtypes = [
        GoFloat64Matrix, # X
        GoFloat64Slice, # y
        GoStringSlice, # X_names
        c_double, # const_max
        c_double, # const_min
        GoString, # eval_metric_name
        GoString, # funcs_string
        GoString, # loss_metric_name
        c_longlong, # max_height
        c_longlong, # min_height
        c_longlong, # n_generations
        c_longlong, # n_populations
        c_longlong, # n_rounds
        c_double, # p_constant
        c_double, # p_full
        c_double, # p_hoist_mutation
        c_double, # p_point_mutation
        c_double, # p_subtree_crossover
        c_double, # p_subtree_mutation
        c_double, # p_terminal
        c_double, # parsimony_coeff
        c_double, # point_mutation_rate
        c_longlong, # population_size
        c_longlong, # seed
        c_longlong, # tuning_n_generations
        c_bool # verbose
    ]

    koza.Fit.restype = c_char_p

    program_bytes = koza.Fit(
        numpy_to_float64_slice(np.transpose(X)),
        numpy_to_float64_slice(y),
        str_list_to_string_slice(X_names),
        const_max,
        const_min,
        GoString(bytes(eval_metric_name, 'utf-8'), len(eval_metric_name)),
        GoString(bytes(funcs_string, 'utf-8'), len(funcs_string)),
        GoString(bytes(loss_metric_name, 'utf-8'), len(loss_metric_name)),
        max_height,
        min_height,
        n_generations,
        n_populations,
        n_rounds,
        p_constant,
        p_full,
        p_hoist_mutation,
        p_point_mutation,
        p_subtree_crossover,
        p_subtree_mutation,
        p_terminal,
        parsimony_coeff,
        point_mutation_rate,
        population_size,
        seed,
        tuning_n_generations,
        verbose
    )

    return program_bytes.decode()
