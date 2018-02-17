import numpy as np
import random
from nqo import world_interface as wi

experience = []
gamma = 0
epsilon = 0
epsilon_min = 0
epsilon_decay = 0


def record_new_experience(state, action, reward, next_state, done):
    experience.append((state, action, reward, next_state, done))


def generate_experience(model, batch_size=16):
    batch = random.sample(experience, batch_size)
    for state, action, reward, next_state, done in batch:
        # if done, make our target reward
        target = reward
        if not done:
            # predict the future discounted reward
            target = reward + gamma * np.amax(model.predict(next_state)[0])
        # make the agent to approximately map
        # the current state to future discounted reward
        # We'll call that target_f
        target_f = model.predict(state)
        target_f[0][action] = target
        # Train the Neural Net with the state and target_f
        model.fit(state, target_f, epochs=1, verbose=0)
    global epsilon
    if epsilon > epsilon_min:
        epsilon *= epsilon_decay


def select_action(model, state, num_actions):
    if np.random.rand() <= epsilon:
        return random.randrange(num_actions)
    act_values = model.predict(state)
    return np.argmax(act_values[0])  # returns action


def life_cycle(model, num_actions, iterations):
    # Iterate the game
    for e in range(iterations):
        # reset state in the beginning of each game
        state = wi.obtain_state()
        # time_t represents each frame of the game
        # Our goal is to keep the pole upright as long as possible until score of 500
        # the more time_t the more score
        for time_t in range(500):
            # turn this on if you want to render
            # env.render()
            # Decide action
            action = select_action(model, state, num_actions)
            # Advance the game to the next frame based on the action.
            # Reward is 1 for every frame the pole survived
            next_state, reward, done, _ = wi.execute_transformation(action)
            # Remember the previous state, action, reward, and done
            record_new_experience(state, action, reward, next_state, done)
            # make next_state the new current state for the next frame.
            state = next_state
            # done becomes True when the game ends
            # ex) The agent drops the pole
            if done:
                # print the score and break out of the loop
                print("episode: {}/{}, score: {}".format(e, iterations, time_t))
                break
        # train the agent with the experience of the episode
        generate_experience(model)
