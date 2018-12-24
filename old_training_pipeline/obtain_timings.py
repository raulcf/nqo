import subprocess
import numpy as np
import time
import pickle

COCKROACH_PATH = "/Users/ra/go/src/github.com/cockroachdb/cockroach/"
DATABASE_NAME = "tpch"


def normalize_execution_time(list_of_execution_time):
    list_of_execution_time = np.asarray(list_of_execution_time)
    max_val = np.max(list_of_execution_time)
    min_val = np.min(list_of_execution_time)
    span = max_val - min_val
    normalized_values = [((x - min_val) / span) for x in list_of_execution_time]
    return normalized_values


def prepare_command(query):
    c1 = "./cockroach sql --insecure --database " \
         "" + str(DATABASE_NAME) + " --execute='" + str(query) + "'"
    return c1


def time_execution_of(query):
    command = prepare_command(query)

    start_execution = time.time()
    process = subprocess.Popen(command, cwd=COCKROACH_PATH, shell=True, stdout=subprocess.PIPE)
    output, error = process.communicate()
    end_execution = time.time()

    execution_time = end_execution - start_execution

    return execution_time


def read_queries(path_to_queries, path_to_encoded_queries, path_to_training_data, debug_cut_queries=None):
    with open(path_to_queries, 'r') as f:
        queries = f.readlines()
    # debug
    if debug_cut_queries:
        queries = queries[:debug_cut_queries]
    num_queries = len(queries)
    print("Found {} queries".format(str(num_queries)))
    query_times = []
    for q in queries:
        q = q.strip()
        q = q + ';'
        print("Executing {}...".format(q))
        runtime = time_execution_of(q)
        print("Took {}".format(str(runtime)))
        print("")
        query_times.append(runtime)
    normalized_query_times = normalize_execution_time(query_times)
    with open(path_to_encoded_queries, 'r') as f:
        encoded_queries = f.readlines()
    if debug_cut_queries:
        encoded_queries = encoded_queries[:debug_cut_queries]
    num_encoded = len(encoded_queries)
    assert(num_encoded == num_queries)
    encoded_queries = [[int(x) for x in eq.split(',')] for eq in encoded_queries]
    training_data = []
    for q, eq, t in zip(queries, encoded_queries, normalized_query_times):
        training_data.append((q.strip(), eq, t))

    print("Storing training data...")
    with open(path_to_training_data, 'wb') as f:
        pickle.dump(training_data, f)
    print("Storing training data...OK")

    return training_data


def __test_call_process():
    # Test bash communication
    print("Testing comm with bash...")
    print("  ")

    example_query = 'select * from tpch.customer limit 1;'

    command = prepare_command(example_query)

    print("command: " + str(command))

    start = time.time()
    process = subprocess.Popen(command, cwd=COCKROACH_PATH, shell=True, stdout=subprocess.DEVNULL)
    output, error = process.communicate()
    end = time.time()

    print(output)
    total_time = end - start
    print("Took: " + str(total_time))


if __name__ == "__main__":
    print("Build training data")

    path_to_query_data = "/Users/ra/dev/old_nqo/raw_query_data/queries.txt"
    path_to_encoded_queries = "/Users/ra/dev/old_nqo/raw_query_data/encoded.txt"
    path_to_training_data = "test"

    training_data = read_queries(path_to_query_data,
                                 path_to_encoded_queries,
                                 path_to_training_data,
                                 debug_cut_queries=None)
    print(str(training_data))


