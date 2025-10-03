import { RJSFSchema } from "@rjsf/utils";
import SchemaEditor from "./SchemaEditor";

const messageSchema = {
  type: "object",
  properties: {
    content: {
      title: "Content",
      type: "string",
    },
    embeds: {
      title: "Embeds",
      type: "array",
      items: {
        type: "object",
        properties: {
          title: {
            title: "Title",
            type: "string",
          },
          description: {
            title: "Description",
            type: "string",
          },
          fields: {
            title: "Fields",
            type: "array",
            items: {
              type: "object",
            },
            properties: {
              name: {
                title: "Name",
                type: "string",
              },
              value: {
                title: "Value",
                type: "string",
              },
            },
          },
        },
      },
    },
  },
} satisfies RJSFSchema;

const messageUISchema = {};

export default function MessageEditor({
  data,
  onChange,
}: {
  data: any;
  onChange: (data: any) => void;
}) {
  return (
    <SchemaEditor
      schema={messageSchema}
      uiSchema={messageUISchema}
      formData={data}
      onChange={onChange}
    />
  );
}
