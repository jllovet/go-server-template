# Overview

As I have been relearning Golang, I have been spending time studying how to build http servers. This repository is an attempt to synthesize lessons and patterns to help make building and managing http servers sensible and able to scale.

Another priority of mine is to build servers that can easily be compared to OpenAPI specs. At some point in the future, I may try to do this through code generation, such as with Oto or one of the many OpenAPI code generation tools. I believe that spec-first development empowers API-first development, which is a valuable way to build products and services in a way that all of the consumers of it - from developers to admins to end users - will be happy and more productive.

At the time of this initial writing, I am refraining from using code generation, because I want to make sure that I'm being intentional about each of the server components I'm using. If you have a recommendation about how to approach integrating a code-generation tool into this project, please open a Github Issue and let me know your thoughts.

# Getting Started

Clone the repository.

Run the following to get a .env file to use.
```shell
cp .env_example .env
```

Using the default values there, you can spin up a server running on `http://localhost:8080` by running

```shell
make build && make run
```

---

# Acknowledgements

I'm heavily indebted to [Mat Ryer](https://github.com/matryer) for the ideas and patterns in this repository. He was a host of [Go Time](https://changelog.com/gotime), has written copiously, and has given a number of talks. That said, I am not proceeding religiously with his suggestions, and there may be aspects of the template that he would agree with. I would love to know what you think about particular decisions, and I'd love to talk about alternatives. As he mentions in this [episode of Go Time](https://www.youtube.com/watch?v=tJ1zvBBkmmY&ab_channel=Changelog), his blog posts and talks are meant to be places to draw ideas from, not things to be followed to the letter. The patterns should only be used if they work for you.

- [How I write HTTP services in Go after 13 years - Grafana Blog](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
- [How Mat writes HTTP services in Go - Changelog: Go Time](https://www.youtube.com/watch?v=tJ1zvBBkmmY&ab_channel=Changelog)
- [GopherCon 2019: How I Write HTTP Web Services after Eight Years - Mat Ryer](https://www.youtube.com/watch?v=rWBSMsLG8po&ab_channel=GopherAcademy)
- [How I build APIs capable of gigantic scale in Go – Mat Ryer](https://www.youtube.com/watch?v=FkPqqakDeRY)

