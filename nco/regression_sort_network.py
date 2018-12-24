from keras.models import Sequential
from keras.layers import Dense, Embedding, LSTM, Input, Dropout
from keras.preprocessing.sequence import pad_sequences
import numpy as np


def build_model(seq_len, vocab_size):
    model = Sequential()
    # model.add(Embedding(input_dim=num_features, output_dim=128, input_length=max_seq_len, name="emb_q"))
    model.add(Dense(256, activation='relu', input_dim=seq_len, input_dtype=np.float32))
    # model.add(Dropout(0.5))
    model.add(Dense(128, activation='relu'))
    # model.add(Embedding(input_dim=seq_len, output_dim=128, input_length=seq_len, name="emb_q"))
    # model.add(LSTM(units=128, name="seq"))
    # base.add(Dropout(0.5))
    # model.add(Dense(128, name='dense', activation='relu'))
    model.add(Dense(1, name='out', activation='linear'))
    # model.add(Dense(1, name='out2', activation='linear'))
    return model


def compile(model):
    model.compile(loss='mean_absolute_error', optimizer='adam', metrics=['mean_absolute_error'])
    # model.compile(loss='mean_squared_error', optimizer='adam', metrics=['mean_squared_error'])
    return model


def read_lines(file_path):
    with open(file_path, 'r') as f:
        lines = f.readlines()
    lines = [l.strip() for l in lines]
    return lines

if __name__ == "__main__":
    print("Sort regression network")

    # read data in
    in_path = "/Users/ra-mit/development/nqo/training_data/sample_training_data.dat"
    # output_path = "/Users/ra-mit/development/nqo/training_data/sample_data"

    lines = read_lines(in_path)

    data = []
    y = []
    for l in lines:
        cost_query = l.split(' ')  # cost <space> <comma_separated query encoding>
        cost = cost_query[0].strip()
        query = cost_query[1]
        toks = query.split(",")
        data.append(toks)
        y.append([float(cost)])

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
    # Y = np.asarray(y, dtype=np.float32)
    Y = [i for i in range(len(X))]
    Y = np.asarray(Y)

    # print(X)

    # split into training test ?
    vocab_size = len(vocab)
    print("Vocab size: " + str(vocab_size))
    model = build_model(max_seq_len, vocab_size)
    model = compile(model)
    model.summary()

    # train the thing
    num_epochs = 1500

    print(Y.shape)
    model.fit(X, Y, epochs=num_epochs, batch_size=32)

    for i, x in enumerate(X):
        res = model.predict(np.asarray([x]))
        print("-> " + str(i))
        print(res)



