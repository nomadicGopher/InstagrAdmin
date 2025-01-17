## Pre-Requisites

* Obviousely, an insagram account.
* **For now**, you must enable the Developer settings associated with that instagram (Meta) account, in order to acquire the User Access Token used in this program. *This allows for your account information to be ethically accessed.*

## Instructions

Preferraby create a `config.json` file in the same directory as the program after it is downloaded that follows the formatting defined in `config-sample.json`. I recommend simply downloading the file from this repository, renaming it, and then setting your preffered values. Spaces should not be included in any of the values and only username & access_token are required, the others will default to false for `include_verified` & the current directory of the program for `output_directory`.

Alternatively, you may run the `UnmutualConnections.exe` with no config file and enter your Intagram `username` & `access_token` when prompted. You may also choose to `include_verified` accounts in the report aswell as set a specific `output_directory` on your PC to save the report to (if you want it to be different from the program).
