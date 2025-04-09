import time
import requests


BASE_URL = "http://localhost:8080"

#Wait initial time
time.sleep(60)


# Test creating a recipedf
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