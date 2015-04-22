# sshfs-mountie

Mountie allows CF applications to mount one or more filesystems provided by the [cf-sshfs](https://github.com/pivotal-cf/cf-sshfs) service.


To use mountie with a CF app that you've already started:

1. Create service instances for your desired filesystems.

   ```
   cf create-service sshfs unlimited my-web-content
   cf create-service sshfs unlimited my-internal-content
   ```
2. Bind the service instances to your application.

   ```
   cf bind-service my-app my-web-content
   cf bind-service my-app my-internal-content
   ```
3. Create a directory called `.profile.d` in your application's root directory.

   ```bash
   cd my-app
   mkdir -p .profile.d
   ```
4. Download the `mount.sh` file from the [`latest release`](https://github.com/pivotal-cf/sshfs-mountie/releases/latest) and place it in your application's `.profile.d/mount.sh`.

5. Push your app.

   ```bash
   cf push my-app
   ```

By default, Mountie will mount all filesystems inside `/home/vcap/filesystems` in the application container.

An application developer may customize this path by modifying the `mount.sh` script inside their `.profile.d` directory (between steps 4 and 5 above).

