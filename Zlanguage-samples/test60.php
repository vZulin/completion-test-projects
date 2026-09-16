<?php

function greeting(string $name): string
{
    return "Hello, {$name}!";
}

echo greeting('PhpStorm');
