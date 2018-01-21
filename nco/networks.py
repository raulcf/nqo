from keras.models import Sequential, Input, Model
from keras.layers import Dense, Embedding
from keras.layers import LSTM
from keras.preprocessing.sequence import pad_sequences
import numpy as np


def lstm_network(input_dim, embedding_dim=128):
    model = Sequential()
    model.add(Embedding(input_dim=input_dim, output_dim=embedding_dim, input_length=50))
    model.add(LSTM(embedding_dim, dropout=0.2, recurrent_dropout=0.2))
    model.add(Dense(1, activation='linear'))

    model.compile(loss='mean_squared_error', optimizer='rmsprop', metrics=['accuracy'])
    return model


def train(model: Sequential, X, Y, epochs, batch_size):
    model.fit(X, Y, epochs=epochs, batch_size=batch_size)
    return model


def train_variable_length_sequences(model, X, Y, epochs):
    for _ in range(epochs):
        for seq, label in zip(X, Y):
            seq = np.asarray([seq])
            model.fit(seq, [label], verbose=1)


if __name__ == "__main__":
    print("Networks")


    def test_variable_sequences():
        input_data_size = 100
        sequence_max_length = 50

        # Create random input sequences. Each sequence will have a different length
        X = []
        for _ in range(input_data_size):
            random_seq_length = np.random.randint(sequence_max_length, size=(1, 1))[0][0] + 1
            xs = [np.random.randint(2, size=(1, 16))[0][0] for _ in range(random_seq_length)]
            X.append(np.asarray(xs))

        Y = np.random.rand(input_data_size, 1)
        model = lstm_network(input_dim=16)
        model.summary()
        train_variable_length_sequences(model, X, Y, epochs=3)

    def test_padding_sequences():
        input_data_size = 100
        int_range = 50

        sequence_max_length = 50

        # Create random input sequences. Each sequence will have a different length
        X = []
        for _ in range(input_data_size):
            random_seq_length = np.random.randint(sequence_max_length, size=(1, 1))[0][0] + 1
            xs = [np.random.randint(int_range, size=(1, 1))[0][0] for _ in range(random_seq_length)]
            X.append(xs)

        X = pad_sequences(X, padding='post', truncating='post', maxlen=sequence_max_length)

        Y = np.random.rand(input_data_size, 1)
        model = lstm_network(input_dim=50, embedding_dim=128)
        model.summary()
        train(model, X, Y, epochs=3, batch_size=10)

    test_padding_sequences()
