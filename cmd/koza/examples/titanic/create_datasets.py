import os

import pandas as pd
from sklearn import model_selection


def munge(data):
    # Sex
    data.drop(['Ticket', 'Name'], inplace=True, axis=1)
    data.Sex.fillna('0', inplace=True)
    data.loc[data.Sex != 'male', 'Sex'] = 0
    data.loc[data.Sex == 'male', 'Sex'] = 1

    # Cabin
    data.Cabin.fillna('0', inplace=True)
    data.loc[data.Cabin.str[0] == 'A', 'Cabin'] = 1
    data.loc[data.Cabin.str[0] == 'B', 'Cabin'] = 2
    data.loc[data.Cabin.str[0] == 'C', 'Cabin'] = 3
    data.loc[data.Cabin.str[0] == 'D', 'Cabin'] = 4
    data.loc[data.Cabin.str[0] == 'E', 'Cabin'] = 5
    data.loc[data.Cabin.str[0] == 'F', 'Cabin'] = 6
    data.loc[data.Cabin.str[0] == 'G', 'Cabin'] = 7
    data.loc[data.Cabin.str[0] == 'T', 'Cabin'] = 8

    # Embarked
    data.loc[data.Embarked == 'C', 'Embarked'] = 1
    data.loc[data.Embarked == 'Q', 'Embarked'] = 2
    data.loc[data.Embarked == 'S', 'Embarked'] = 3
    data.Embarked.fillna(0, inplace=True)
    data.fillna(-1, inplace=True)

    return data

if __name__ == '__main__':

    train = munge(pd.read_csv('kaggle/train.csv'))
    test = munge(pd.read_csv('kaggle/test.csv'))
    train, val = model_selection.train_test_split(train, test_size=0.2, random_state=42)

    here = os.path.dirname(os.path.realpath(__file__))
    train.to_csv(os.path.join(here, 'train.csv'), index=False)
    val.to_csv(os.path.join(here, 'val.csv'), index=False)
    test.to_csv(os.path.join(here, 'test.csv'), index=False)
