import pandas as pd
from sklearn.utils import check_random_state


if __name__ == '__main__':

    rng = check_random_state(0)

    # Training samples
    X_train = rng.uniform(-1, 1, 1000).reshape(500, 2)
    y_train = X_train[:, 0]**2 - X_train[:, 1]**2 + X_train[:, 1] - 1
    pd.dataset({
        'x0': X_train[:, 0],
        'x1': X_train[:, 1],
        'y': y_train
    }).to_csv('train.csv', index=False)

    # Testing samples
    X_test = rng.uniform(-1, 1, 100).reshape(50, 2)
    y_test = X_test[:, 0]**2 - X_test[:, 1]**2 + X_test[:, 1] - 1
    pd.dataset({
        'x0': X_test[:, 0],
        'x1': X_test[:, 1],
        'y': y_test
    }).to_csv('test.csv', index=False)
