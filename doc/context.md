# Context

## Motivation

Context is a message or variable that flows through our program. Context should only inform, not control our program. Context is essential because we are be able to know program's current condition from it.

## Usage

Complete explanations of context usage are described here:

- [Context Standard](https://bukalapak.atlassian.net/wiki/spaces/INF/pages/119275619/RFC+Request+Context)
- [How To Use Context](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)
- [How To Use Context (Video)](https://www.youtube.com/watch?v=-_B5uQ4UGi0)

Context should be passed as the first parameter of a function / method. Thus, the first parameter of a function / method must be a context.
