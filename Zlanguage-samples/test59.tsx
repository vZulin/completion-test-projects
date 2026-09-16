type GreetingProps = {
  name: string;
};

export function Greeting({ name }: GreetingProps): JSX.Element {
  return <h1>Hello, {name}!</h1>;
}

export default function App(): JSX.Element {
  return <Greeting name="TypeScript" />;
}
