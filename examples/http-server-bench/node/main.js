import express from "express";
import fs from "fs/promises";

const app = express();

app.get("/user/:id", async (req, res) => {
  const { id } = req.params;
  const usersFile = await fs.readFile("./main.json");
  const users = JSON.parse(usersFile.toString());
  const user = users.find((u) => u.id === Number(id));
  res.json(user);
});

console.log("Server running on port 3001");
app.listen(3001);
