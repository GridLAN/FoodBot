# FoodBot
Feeling hungry? FoodBot is a Discord bot that helps you satisfy your primal needs.. By helping you select what to eat, randomly!

Add FoodBot to your Discord server [here](https://discord.com/api/oauth2/authorize?client_id=921912124765253653&permissions=377957124096&scope=bot).

## Development:
1. Compile and run the project.

    ```
    TOKEN=abc123 go run main.go
    ```

2. Alternatively, build and run the project inside of a container.

    ```
    docker build -t foodbot . && docker run --env-file .env foodbot
    ```