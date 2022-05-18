# 💥 go-bloc
![lint_and_test_workflow](https://github.com/hijgo/go-bloc/actions/workflows/lint_and_test.yaml/badge.svg)<br/>
This package provides an adaption of the BloC Design Pattern for go.<br/>
You can use the library to structure your code easily and increase the maintainability of your project at the same time.<br/>

## What is the BloC Design Pattern? 😅
Generally speaking, BLoC (Business Logic Component) is a design pattern enabling developers to efficiently and conveniently manage state across their apps without a tight coupling between the presentation (view) and the logic. It also aims for reusability of the same logic across multiple widgets. It was first mentioned by Google at the Google I/O in 2018.
 -- <cite>[flutterclutter.dev](https://www.flutterclutter.dev/flutter/basics/what-is-the-bloc-pattern/2021/2084/)</cite>
 <br/>
 <br/>
 For further information about the BloC Design Pattern I recommend [mitrais.com](https://www.mitrais.com/news-updates/getting-started-with-flutter-bloc-pattern/)
 
## How can it be used in go? 
Although the BloC Design Pattern was introduced for the development of mobile apps using Flutter, it can be beneficial outside its origins.
<br/><br/>
BloC can be useful in go the same way as it can be inside Flutter. Whenever your project is presenting something, in some way or another, you should think about using the BloC Design Pattern. Because without a clear and easy to understand design-pattern, like BloC, it can be hard to understand and keep track of the strings connecting your presenting layer and data layer.<br/><br/>
When the presenting layer and the data layer are separated, with a BloC as a middleman, your project structure will become much clearer. Therefore, your code will become more readable and easier to understand too. As extra plus point testing will also become much easier, as you can produce any state through feeding events to the BloC.

## Want to start your BloC-Journey?
If you want to start with the go-bloc library, please check out the [Example](https://github.com/hijgo/go-bloc-examples) repository to get a head start.
