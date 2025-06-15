import express from "express";
import fs from "fs/promises";

const app = express();

app.get("/user/:id", async (req, res) => {
  const { id } = req.params;
  const users = await fs
    .readFile("./main.json")
    .then((b) => b.toString())
    .then((s) => JSON.parse(s));

  const user = users.find((u) => u.id === Number(id));
  res.json(user);
});

console.log("Server running on port 3000");
app.listen(3001);
