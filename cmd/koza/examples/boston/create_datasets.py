import os

import pandas as pd
from sklearn import datasets
from sklearn import model_selection


if __name__ == '__main__':

    X, y = datasets.load_boston(return_X_y=True)

    df = pd.DataFrame(X)
    df['y'] = y

    train, test = model_selection.train_test_split(df, test_size=0.33, random_state=10)

    here = os.path.dirname(os.path.realpath(__file__))
    train.to_csv(os.path.join(here, 'train.csv'), index=False)
    test.to_csv(os.path.join(here, 'test.csv'), index=False)
