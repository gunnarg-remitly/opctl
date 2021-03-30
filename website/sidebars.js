module.exports = {
  docs: [
    "introduction",
    {
      type: "category",
      label: "Setup",
      items: [
        "setup/bare-metal",
        "setup/azure-pipelines",
        "setup/docker",
        "setup/github",
        "setup/gitlab",
        "setup/kubernetes",
        "setup/travis"
      ]
    },
    {
      type: "category",
      label: "Training",
      items: [
        "training/hello-world",
        "training/inputs-outputs",
        {
          type: "category",
          label: "Portable Dev",
          items: [
            "run-a-go-service",
            "run-a-react-app"
          ]
        },
        {
          type: "category",
          label: "Containers",
          items: [
            "training/containers/how-do-i-communicate-with-an-opctl-container",
            "training/containers/how-do-i-get-opctl-containers-to-communicate",
            "training/containers/how-do-i-run-a-container"
          ]
        }
      ]
    },
    {
      type: "category",
      label: "Reference",
      items: [
        {
          type: "category",
          label: "Opspec",
          items: [
            "reference/opspec/index",
            {
              type: "category",
              label: "op.yml",
              items: [
                "reference/opspec/op.yml/index",
                {
                  type: "category",
                  label: "Call",
                  items: [
                    "reference/opspec/op.yml/call/index",
                    {
                      type: "category",
                      label: "Container call",
                      items: [
                        "reference/opspec/op.yml/call/container/index",
                        "reference/opspec/op.yml/call/container/image",
                      ]
                    },
                    "reference/opspec/op.yml/call/op",
                    "reference/opspec/op.yml/call/parallel-loop",
                    "reference/opspec/op.yml/call/serial-loop",
                  ]
                },
                {
                  type: "category",
                  label: "Parameter",
                  items: [
                    "reference/opspec/op.yml/parameter/index",
                    "reference/opspec/op.yml/parameter/array",
                    "reference/opspec/op.yml/parameter/boolean",
                    "reference/opspec/op.yml/parameter/dir",
                    "reference/opspec/op.yml/parameter/file",
                    "reference/opspec/op.yml/parameter/number",
                    "reference/opspec/op.yml/parameter/object",
                    "reference/opspec/op.yml/parameter/socket",
                    "reference/opspec/op.yml/parameter/string",
                  ]
                },
                "reference/opspec/op.yml/identifier",
                "reference/opspec/op.yml/initializer",
                "reference/opspec/op.yml/loop-vars",
                "reference/opspec/op.yml/markdown",
                "reference/opspec/op.yml/predicate",
                "reference/opspec/op.yml/pull-creds",
                "reference/opspec/op.yml/rangeable-value",
                "reference/opspec/op.yml/variable-reference"
              ]
            },
            {
              type: "category",
              label: "Types",
              items: [
                "reference/opspec/types/array",
                "reference/opspec/types/boolean",
                "reference/opspec/types/dir",
                "reference/opspec/types/file",
                "reference/opspec/types/number",
                "reference/opspec/types/object",
                "reference/opspec/types/socket",
                "reference/opspec/types/string",
              ]
            }
          ]
        },
        {
          type: "category",
          label: "CLI",
          items: [
            "reference/cli/global-options",
            {
              type: "category",
              label: "auth",
              items: [
                "reference/cli/auth/index",
                "reference/cli/auth/add",
              ]
            },
            "reference/cli/ls",
            {
              type: "category",
              label: "op",
              items: [
                "reference/cli/op/index",
                "reference/cli/op/create",
                "reference/cli/op/install",
                "reference/cli/op/validate",
              ]
            },
            "reference/cli/run",
            "reference/cli/self-update",
          ]
        },
        {
          type: "link",
          label: "ReST API",
          href: "https://petstore.swagger.io/?url=https://raw.githubusercontent.com/opctl/opctl/main/api/openapi.yaml"
        },
        {
          type: "category",
          label: "SDKs",
          items: [
            {
              type: "link",
              label: "Go",
              href: "https://github.com/opctl/opctl/tree/main/sdks/go"
            },
            {
              type: "link",
              label: "Js",
              href: "https://github.com/opctl/opctl/tree/main/sdks/js"
            },
            {
              type: "link",
              label: "React",
              href: "https://github.com/opctl/opctl/tree/main/sdks/react"
            }
          ]
        }
      ]
    }
  ]
};
