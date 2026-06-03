<?php
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Slim\Factory\AppFactory;

require __DIR__ . '/vendor/autoload.php';

$app = AppFactory::create();

$app->get('/', function (Request $request, Response $response) {
    $all_envs = getenv();
    $formatted = [];

    foreach ($all_envs as $key => $value) {
        $formatted[] = "$key=$value"; // Format as KEY=VALUE
    }
    // Implode using a comma or a newline
    $envString = implode('; ', $formatted);
    $response->getBody()->write($envString);
    // $result = random_int(1,6);
    // $response->getBody()->write(strval($result));
    return $response;
});

$app->run();
