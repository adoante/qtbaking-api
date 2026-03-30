import json
from os import listdir
from os.path import join
import psycopg2

# connect to db
# get values from the .env
conn = psycopg2.connect(
    host="localhost",
    database="qtbaking",
    user="postgres",
    password="https://www.youtube.com/watch?v=XdNd0IwiGEo",
    port=5432
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
            """, (recipe["slug"], recipe["title"], recipe["video_url"], recipe["created_at"]))

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

        cur.execute("""
        INSERT INTO recipes (vod_id, thumbnail, temp_fahrenheit, temp_celsius)
        VALUES (%s, %s, %s, %s)
        RETURNING id
        """, (vod_id, thumbnail, temp_fahrenheit, temp_celsius))

        recipe_id = cur.fetchone()[0]
        recipes[recipe["title"]] = recipe_id

        f.close()

# components

# ingredients

# tools

# notes

# tags

# close connection !
conn.commit()
cur.close()
conn.close()
