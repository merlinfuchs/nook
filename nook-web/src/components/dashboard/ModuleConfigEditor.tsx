import { RJSFSchema, RJSFValidationError, UiSchema } from "@rjsf/utils";
import SchemaEditor from "../SchemaEditor";

export default function ModuleConfigEditor({
  schema,
  uiSchema,
  formData,
  onChange,
}: {
  schema: RJSFSchema;
  uiSchema?: UiSchema;
  formData?: any;
  onChange: (data: any, errors: RJSFValidationError[]) => void;
}) {
  return (
    <SchemaEditor
      schema={schema}
      uiSchema={uiSchema}
      formData={formData}
      onChange={onChange}
    />
  );
}
