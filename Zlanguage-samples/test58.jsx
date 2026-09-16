import React from "react";

export function Greeting({ name }) {
  return <h1>Hello, {name}!</h1>;
}

export default function App() {
  return <Greeting name="WebStorm" />;
}
