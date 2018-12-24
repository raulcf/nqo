from keras import Sequential
from keras.layers import Embedding, Dense, LSTM


def define_model(input_dim, action_dim, embedding_dim=128):

    model = Sequential()
    model.add(Embedding(input_dim=input_dim, output_dim=embedding_dim, input_length=50))
    model.add(LSTM(embedding_dim, dropout=0.2, recurrent_dropout=0.2))
    model.add(Dense(action_dim, activation='softmax'))

    model.compile(loss='categorical_crossentropy', optimizer='rmsprop', metrics=['accuracy'])
    return model

