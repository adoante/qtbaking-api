import json
from os import listdir
from os.path import join
import psycopg2
from psycopg2.sql import NULL

# connect to db
# get values from the .env
conn = psycopg2.connect(
    host="localhost",
    database="qtbaking",
    user="postgres",
    password="https://www.youtube.com/watch?v=XdNd0IwiGEo",
    port=4321
)

cur = conn.cursor()

# load the recipes json
recipes_json = [f for f in listdir("./recipes-json") if f.endswith(".json")]


# add vods
# 1 vod can have many recipes
# many recipes can have 1 vod

vods = {}

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        if recipe["video_url"] and recipe["video_url"] not in vods:
            cur.execute("""
            INSERT INTO vods (slug, title, video_url, created_at)
            VALUES (%s, %s, %s, %s)
            RETURNING id
            """, (recipe["slug"], recipe["vod_title"], recipe["video_url"], recipe["created_at"]))

            vod_id = cur.fetchone()[0]
            vods[recipe["video_url"]] = vod_id

            print("Added ", file_name, vods[recipe["video_url"]])
        else:
            print("*", file_name, recipe["video_url"])

        f.close()

# add recipe(s)
# recipes may or may not have vods

recipes = {}

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        vod_id = vods[recipe["video_url"]] if recipe["video_url"] else None
        thumbnail = recipe["thumbnail"] if recipe["thumbnail"] else None
        temp_fahrenheit = recipe["temp_fahrenheit"] if recipe["temp_fahrenheit"] else None
        temp_celsius = recipe["temp_celsius"] if recipe["temp_celsius"] else None
        title = recipe["title"]

        cur.execute("""
        INSERT INTO recipes (vod_id, thumbnail, title, temp_fahrenheit, temp_celsius)
        VALUES (%s, %s, %s, %s, %s)
        RETURNING id
        """, (vod_id, thumbnail, title, temp_fahrenheit, temp_celsius))

        recipe_id = cur.fetchone()[0]
        recipes[recipe["title"]] = recipe_id

        f.close()

# components

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        recipe_id = recipes[recipe["title"]]

        # loop through components array
        for component in recipe["components"]:
            name = component["name"]

            cur.execute("""
            INSERT INTO components (recipe_id, name)
            VALUES (%s, %s)
            RETURNING id
            """, (recipe_id, name))

            component_id = cur.fetchone()[0]

            # ingredients
            for ingredient in component["ingredients"]:
                name = ingredient["name"]
                quantity = ingredient["quantity"]
                unit = ingredient["unit"]
                metric_quantity = ingredient["metric_quantity"] if ingredient["metric_quantity"] else None
                metric_unit = ingredient["metric_unit"] if ingredient["metric_unit"] else None
                optional = ingredient["optional"]
                notes = ", ".join(
                    ingredient["notes"]) if ingredient["notes"][0] != "" else None

                cur.execute("""
                INSERT INTO ingredients (component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                """, (component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes))


# tools

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        recipe_id = recipes[recipe["title"]]

        for tool in recipe["tools"]:
            name = tool["name"]
            optional = tool["optional"]

            cur.execute("""
            INSERT INTO tools (recipe_id, name, optional)
            VALUES (%s, %s, %s)
            """, (recipe_id, name, optional))


# notes

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        recipe_id = recipes[recipe["title"]]

        for note in recipe["notes"]:
            if note != "":
                cur.execute("""
                INSERT INTO notes (recipe_id, note)
                VALUES (%s, %s)
                """, (recipe_id, note))

# tags

for file_name in recipes_json:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        recipe_id = recipes[recipe["title"]]

        for tag in recipe["tags"]:

            cur.execute("""
            INSERT INTO tags (recipe_id, tag)
            VALUES (%s, %s)
            """, (recipe_id, tag))

# close connection !
conn.commit()
cur.close()
conn.close()
