import itertools
import random


TRAINING_DATA_SEPARATOR = "-%%%-"

def generate_training_samples(ordered_plans):
    training_samples = []
    for f, s in itertools.combinations(ordered_plans, 2):
        training_samples.append((s, f, 1))
    for s, f in itertools.combinations(reversed(ordered_plans), 2):
        training_samples.append((f, s, 0))
    random.shuffle(training_samples)
    return training_samples


def read_lines(file_path):
    with open(file_path, 'r') as f:
        lines = f.readlines()
    lines = [l.strip() for l in lines]
    return lines


def write_training_data(output_path, samples):
    with open(output_path, 'w') as f:
        for a, b, l in samples:
            line = str(a) + TRAINING_DATA_SEPARATOR + str(b) + TRAINING_DATA_SEPARATOR + str(l) + '\n'
            f.write(line)
    print("Written to: " + str(output_path))


if __name__ == "__main__":
   print("Prepare training data")

   # test = [1,2,3,4,5,6]
   #
   # result = generate_training_samples(test)
   # print(len(result))

   input_path = "/Users/ra-mit/development/nqo/training_data/sample_data"
   output_path = "/Users/ra-mit/development/nqo/training_data/sample_training_data.dat"

   lines = read_lines(input_path)
   training_data = generate_training_samples(lines)
   write_training_data(output_path, training_data)



