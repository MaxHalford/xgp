import pandas as pd
from sklearn import datasets
from sklearn import model_selection


if __name__ == '__main__':

    X, y = datasets.load_boston(return_X_y=True)

    df = pd.DataFrame(X)
    df['y'] = y

    train, test = model_selection.train_test_split(df, test_size=0.33, random_state=42)

    train.to_csv('train.csv', index=False)
    test.to_csv('test.csv', index=False)
