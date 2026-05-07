def moving_average(prices, window):
    if not prices or len(prices) < window:
        return None
    return sum(prices[-window:]) / window


def change_percent(prices):
    if len(prices) < 2 or prices[-2] == 0:
        return None
    return (prices[-1] - prices[-2]) / prices[-2] * 100


def volatility(prices):
    if len(prices) < 2:
        return None
    changes = []
    for prev, cur in zip(prices, prices[1:]):
        if prev != 0:
            changes.append((cur - prev) / prev * 100)
    if not changes:
        return None
    mean = sum(changes) / len(changes)
    variance = sum((item - mean) ** 2 for item in changes) / len(changes)
    return variance ** 0.5
