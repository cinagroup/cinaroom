(how-to-guides-customise-cinaclaw-authenticate-users-with-the-cinaclaw-service)=
# Authenticate users with the cinaclaw service

> See also: [`authenticate`](reference-command-line-interface-authenticate), [local.passphrase](reference-settings-local-passphrase), [Service](explanation-service)

cinaclaw requires users to be authenticated with the service before allowing commands to complete. The installing user is automatically authenticated.

## Setting the passphrase

The administrator needs to set a passphrase for users to authenticate with the cinaclaw service. The user setting the passphrase will need to already be authenticated.

There are two ways to proceed:

* Set the passphrase with an echoless interactive entry, where the passphrase is hidden from view:

   ```{code-block} text
   cinaclaw set local.passphrase
   ```

   The system will then prompt you to enter a passphrase:

   ```{code-block} text
   Please enter passphrase:
   Please re-enter passphrase:
   ```

* Set the passphrase in the command line, where the passphrase is visible:

   ```{code-block} text
   cinaclaw set local.passphrase=foo
   ```

## Authenticating the user

A user that is not authorized to connect to the cinaclaw service will fail when running `cinaclaw` commands. An error will be displayed when this happens.

For example, if you try running the `cinaclaw list` command:

```{code-block} text
list failed: The user is not authenticated with the cinaclaw service.

Please authenticate before proceeding (e.g. via 'cinaclaw authenticate'). Note that you first need an authenticated user to set and provide you with a trusted passphrase (e.g. via 'cinaclaw set local.passphrase').
```

At this time, the user will need to provide the previously set passphrase. This can be accomplished in two ways:

* Authenticate with an echoless interactive entry, where the passphrase is hidden from view:

    ```{code-block} text
    cinaclaw authenticate
    ```

    The system will prompt you to enter the passphrase:

     ```{code-block} text
    Please enter passphrase:
    ```

* Authenticate in the command line, where the passphrase is visible:

   ```{code-block} text
   cinaclaw authenticate foo
   ```

## Troubleshooting

Here you can find solutions and workarounds for common issues that may arise.

### The user cannot be authorized and the passphrase cannot be set

It is possible that another user that is privileged to connect to the cinaclaw socket will
connect first and make it seemingly impossible to set the `local.passphrase` and also `authorize`
the user with the service. This usually occurs when cinaclaw is installed as `root/admin` but
the user is run as another user, or vice versa.

If this is the case, you will see something like the following when you run:

* `cinaclaw list`

  ```{code-block} text
  list failed: The user is not authenticated with the cinaclaw service.

  Please authenticate before proceeding (e.g. via 'cinaclaw authenticate'). Note that you first need an authenticated user to set and provide you with a trusted passphrase (e.g. via 'cinaclaw set local.passphrase').
  ```

* `cinaclaw authenticate`

  ```{code-block} text
  Please enter passphrase:
  authenticate failed: No passphrase is set.

  Please ask an authenticated user to set one and provide it to you. They can achieve so with 'cinaclaw set local.passphrase'. Note that only the user who installs cinaclaw is automatically authenticated.
  ```

* `cinaclaw set local.passphrase`

  ```{code-block} text
  Please enter passphrase:
  Please re-enter passphrase:
  set failed: The user is not authenticated with the cinaclaw service.

  Please authenticate before proceeding (e.g. via 'cinaclaw authenticate'). Note that you first need an authenticated user to set and provide you with a trusted passphrase (e.g. via 'cinaclaw set local.passphrase').
  ```

This may not even work when using `sudo`.

The following workaround should help get out of this situation:

```bash
cat ~/snap/cinaclaw/current/data/cinaclaw-client-certificate/cinaclaw_cert.pem | sudo tee -a /var/snap/cinaclaw/common/data/cinaclawd/authenticated-certs/cinaclaw_client_certs.pem > /dev/null

snap restart cinaclaw
```

You may need `sudo` with this last command: `sudo snap restart cinaclaw`.

At this point, your user should be authenticated with the cinaclaw service.
