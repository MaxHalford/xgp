import os

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

    here = os.path.dirname(os.path.realpath(__file__))
    train.to_csv(os.path.join(here, 'train.csv'), index=False)
    test.to_csv(os.path.join(here, 'test.csv'), index=False)
