from keras.models import Sequential, Model


def build_model(input_dim, num_features):
    base = Sequential()
    # SEQUENCE
    # for pos model 128 in old 3 next layers is good
    base.add(Embedding(num_features, output_dim=128, name="emb_q"))  # 64
    base.add(LSTM(units=128, return_sequences=True, name="seq"))  # 32
    # base.add(Dropout(0.5))
    base.add(LSTM(units=128, return_sequences=False, name='seq2'))
    # base.add(LSTM(units=64, return_sequences=True, name="seq3"))
    # base.add(LSTM(units=64, return_sequences=True, name="seq4"))
    # base.add(LSTM(units=32, return_sequences=False, name="seq5"))
    # base.add(Dropout(0.5))

    # NORMAL LINEAR
    # base.add(Dense(1024, activation='relu'))
    # base.add(Dropout(0.2))
    # base.add(Dense(128, activation='relu'))
    # # base.add(Dropout(0.2))
    # base.add(Dense(64, activation='relu'))

    input_a = Input(shape=(input_dim,))
    input_b = Input(shape=(input_dim,))

    siamese_a = base(input_a)
    siamese_b = base(input_b)

    # distance = Lambda(euclidean_distance, output_shape=eucl_dist_output_shape)([siamese_a, siamese_b])

    # merged = merge.Merge([siamese_a, siamese_b], mode=euclidean_distance, output_shape=eucl_dist_output_shape)
    # output = Dense(1, activation='sigmoid')(distance)

    model = Model([input_a, input_b], distance)

    return model