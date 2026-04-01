(explanation-index)=
# Explanation

The following guides provide conceptual context and clarification on the key topics related to using and configuring cinaclaw.

## Architecture

These topics cover the foundations of how cinaclaw operates on your machine, providing a high-level overview of its structure and components.

- [Reference architecture](explanation-reference-architecture): A high-level overview of how cinaclaw is structured, including its clients, daemon, storage, instances, and networking.
- [Host](explanation-host)
- [Platform](explanation-platform)
- [Service](explanation-service)
- [Driver](explanation-driver)



## Instances

These guides explain the lifecycle, identity, and resources of the virtual machines you create.

- [Instance](explanation-instance)
- [Image](explanation-image)
- [Settings keys and values](explanation-settings-keys-values)
- [Blueprint (removed)](explanation-blueprint)

## Using cinaclaw

Concepts related to interacting and extending the functionality of your instances.

- [cinaclaw exec and shells](explanation-cinaclaw-exec-and-shells)
- [Mount](explanation-mount)
- [Alias](explanation-alias)
- [Snapshot](explanation-snapshot)

In cinaclaw, an **alias** is a shortcut for a command that runs inside a given instance.

## Security and performance

Technical background on protecting your environment and ensuring it runs efficiently.

- [About security](explanation-about-security)
- [Authentication](explanation-authentication)
- [ID mapping](explanation-id-mapping)
- [About performance](explanation-about-performance)

---

## Glossary

(explanation-alias)=
### Alias

``` {seealso}
See also: [`alias`](/reference/command-line-interface/alias), [How to use command aliases](/how-to-guides/manage-instances/use-instance-command-aliases).
```
In cinaclaw, an **alias** is a shortcut for a command that runs inside a given instance.


(explanation-host)=
### Host

In cinaclaw, **host** refers the actual physical machine on which cinaclaw is running.


```{toctree}
:titlesonly:
:maxdepth: 2
:hidden:

reference-architecture
platform
service
driver
instance
image
settings-keys-values
blueprint
cinaclaw-exec-and-shells
mount
snapshot
about-security
authentication
id-mapping
about-performance
```
