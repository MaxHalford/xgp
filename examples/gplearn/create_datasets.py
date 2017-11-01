import numpy as np
import pandas as pd
from sklearn import model_selection


if __name__ == '__main__':

    x0 = np.arange(-1, 1, 0.1)
    x1 = np.arange(-1, 1, 0.1)
    y = x0 ** 2 - x1 ** 2 + x1 - 1

    df = pd.DataFrame({
        'x0': x0,
        'x1': x1,
        'y': y
    })

    train, test = model_selection.train_test_split(df, test_size=0.2, random_state=42)

    train.to_csv('train.csv', index=False)
    test.to_csv('test.csv', index=False)
