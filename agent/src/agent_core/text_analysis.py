RISK_KEYWORDS = {
    "loss": "loss",
    "investigation": "investigation",
    "lawsuit": "lawsuit",
    "downgrade": "downgrade",
    "亏损": "loss",
    "调查": "investigation",
    "诉讼": "lawsuit",
    "下调": "downgrade",
}


def extract_risks(texts):
    found = []
    joined = "\n".join(texts).lower()
    for keyword, marker in RISK_KEYWORDS.items():
        if keyword in joined and marker not in found:
            found.append(marker)
    return found
