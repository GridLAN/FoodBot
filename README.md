# FoodBot
Feeling hungry? FoodBot is a Discord bot that helps you satisfy your primal needs.. By helping you select what to eat, randomly!

Data from: https://www.themealdb.com/

## Development:
1. Build and Run to compile the project. This can be slow on the first run.

    ```
    TOKEN=... go run main.go
    ```

2. Alternatively, build and run the project inside of a container.

    ```
    docker build -t kuthero/foodbot . && docker run --env-file .env kuthero/foodbot
    ```