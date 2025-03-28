import requests
import time

BASE_URL = "http://localhost:8080"

# Retry a few times if needed
for i in range(10):
    try:
        r = requests.get(BASE_URL + "/health")
        if r.status_code == 200:
            break
    except Exception:
        time.sleep(2)
else:
    raise Exception("Service did not become healthy")

# Test creating a recipe
response = requests.post(
    BASE_URL + "/recipes",
    json={
        "name": "Spaghetti",
        "description": "Tomato and garlic",
        "ingredients": ["noodles", "tomato", "garlic"]
    }
)
assert response.status_code == 200

# Test fetching all recipes
response = requests.get(BASE_URL + "/recipes")
assert response.status_code == 200
assert "Spaghetti" in response.text

print("Tests passed")