from keras.models import Sequential
from keras.layers import Dense, Embedding, LSTM, Input
from keras.preprocessing.sequence import pad_sequences
import numpy as np


def build_model(seq_len, vocab_size):
    model = Sequential()
    # model.add(Embedding(input_dim=num_features, output_dim=128, input_length=max_seq_len, name="emb_q"))
    model.add(Dense(16, activation='relu', input_dim=seq_len, input_dtype=np.float32))
    # model.add(Dense(128, activation='relu'))
    # model.add(Embedding(input_dim=seq_len, output_dim=128, input_length=seq_len, name="emb_q"))
    # model.add(LSTM(units=128, name="seq"))
    # base.add(Dropout(0.5))
    # model.add(Dense(128, name='dense', activation='relu'))
    model.add(Dense(1, name='out', activation='linear'))
    return model


def compile(model):
    model.compile(loss='mean_absolute_error', optimizer='sgd', metrics=['mean_absolute_error'])
    return model

if __name__ == "__main__":
    print("Sort regression network")

    # read data in
    in_path = "/Users/ra-mit/development/nqo/training_data/sample_data"
    with open(in_path, 'r') as f:
        lines = f.readlines()
    lines = [l.strip() for l in lines]
    data = []
    for l in lines:
        toks = l.split(",")
        data.append(toks[:3])

    # normalize data into sequences
    max_seq_len = 0
    vocab_size = 0
    vocab = set()
    for d in data:
        if len(d) > max_seq_len:
            max_seq_len = len(d)
        for el in d:
            vocab.add(el)
    X = pad_sequences(data, padding='post', truncating='post', maxlen=max_seq_len, dtype=np.float32)
    y = [[12.0], [10.0], [8.0], [6.0], [4.0], [2.0]]
    Y = np.asarray(y, dtype=np.float32)

    # print(X)

    # split into training test ?
    vocab_size = len(vocab)
    print("Vocab size: " + str(vocab_size))
    model = build_model(max_seq_len, vocab_size)
    model = compile(model)
    model.summary()

    # train the thing
    num_epochs = 50

    print(Y.shape)
    model.fit(X, Y, epochs=num_epochs, batch_size=1)

    for i, x in enumerate(X):
        res = model.predict(np.asarray([x]))
        print("-> " + str(i))
        print(res)



