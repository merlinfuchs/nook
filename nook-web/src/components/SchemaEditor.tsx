import { useGuildChannels, useGuildRoles } from "@/lib/hooks/api";
import { ConfigUISelectValue } from "@/lib/types/module.gen";
import Form from "@rjsf/shadcn";
import {
  ArrayFieldTemplateProps,
  FieldErrorProps,
  FieldProps,
  FieldTemplateProps,
  ObjectFieldTemplateProps,
  RegistryWidgetsType,
  RJSFSchema,
  RJSFValidationError,
  UiSchema,
  WidgetProps,
} from "@rjsf/utils";
import validator from "@rjsf/validator-ajv8";
import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChevronDownIcon,
  TrashIcon,
} from "lucide-react";
import { useMemo } from "react";
import MessageEditor from "./MessageEditor";
import { Button } from "./ui/button";
import { Card } from "./ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Input } from "./ui/input";
import { ScrollArea } from "./ui/scroll-area";
import { Select, SelectContent, SelectItem, SelectTrigger } from "./ui/select";
import { Switch } from "./ui/switch";
import MessageFormatEditor from "./MessageFormatEditor";

const SelectWidget = function (props: WidgetProps) {
  const allowMultiple = props.uiSchema?.["ui:allow_multiple"] as boolean;
  const selectValues = props.uiSchema?.[
    "ui:select_values"
  ] as ConfigUISelectValue[];
  if (!selectValues) {
    return null;
  }

  if (allowMultiple) {
    const value = props.value as string[];

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline">
            {props.value?.length > 0
              ? `${props.value.length} selected`
              : "Select a value"}
            <ChevronDownIcon />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          {selectValues.map((item) => (
            <DropdownMenuCheckboxItem
              key={item.value}
              checked={props.value?.includes(item.value)}
              onCheckedChange={(checked) =>
                props.onChange(
                  checked
                    ? [...value, item.value]
                    : value?.filter((i) => i !== item.value)
                )
              }
            >
              {item.label}
            </DropdownMenuCheckboxItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  } else {
    return (
      <Select
        value={props.value ?? ""}
        onValueChange={(value) => props.onChange(value)}
      >
        <SelectTrigger className="px-5">
          {props.value ? props.value : "Select a value"}
        </SelectTrigger>
        <SelectContent>
          {selectValues.map((value) => (
            <SelectItem key={value.value} value={value.value}>
              {value.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  }
};

const ChannelSelectWidget = function (props: WidgetProps) {
  const channels = useGuildChannels();

  const filteredChannels = useMemo(() => {
    if (!props.uiSchema) {
      return channels;
    }

    const channelTypes = props.uiSchema?.["ui:channel_types"];
    if (!channelTypes) {
      return channels;
    }

    return channels?.filter((channel) => channelTypes.includes(channel.type));
  }, [channels, props.uiSchema]);

  const selectedChannel = useMemo(
    () => channels?.find((channel) => channel.id === props.value),
    [channels, props.value]
  );

  return (
    <Select
      value={props.value ?? ""}
      onValueChange={(value) => props.onChange(value)}
    >
      <SelectTrigger className="px-5">
        {props.value ? `#${selectedChannel?.name}` : "Select a channel"}
      </SelectTrigger>
      <SelectContent>
        {filteredChannels?.map((channel) => (
          <SelectItem key={channel.id} value={channel.id}>
            #{channel.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

const RoleSelectWidget = function (props: WidgetProps) {
  const allowMultiple = props.uiSchema?.["ui:allow_multiple"] as boolean;
  const roles = useGuildRoles();

  if (allowMultiple) {
    const value = props.value as string[];

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline">
            {props.value?.length > 0
              ? `${props.value.length} role${
                  props.value.length > 1 ? "s" : ""
                } selected`
              : "Select a role"}
            <ChevronDownIcon />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          {roles?.map((role) => (
            <DropdownMenuCheckboxItem
              key={role.id}
              checked={props.value?.includes(role.id)}
              style={{ color: role.color }}
              onCheckedChange={(checked) =>
                props.onChange(
                  checked
                    ? [...value, role.id]
                    : value?.filter((i) => i !== role.id)
                )
              }
            >
              {role.name}
            </DropdownMenuCheckboxItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  } else {
    const selectedRole = roles?.find((role) => role.id === props.value);

    return (
      <Select
        value={props.value ?? ""}
        onValueChange={(value) => props.onChange(value)}
      >
        <SelectTrigger className="px-5" style={{ color: selectedRole?.color }}>
          {props.value ? selectedRole?.name : "Select a role"}
        </SelectTrigger>
        <SelectContent>
          {roles?.map((role) => (
            <SelectItem
              key={role.id}
              value={role.id}
              style={{ color: role.color }}
            >
              {role.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  }
};

const MessageWidget = function (props: WidgetProps) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Edit Message</Button>
      </DialogTrigger>
      <DialogContent className="w-full sm:max-w-2xl xl:max-w-3xl pr-0">
        <ScrollArea className="h-[80vh] pr-5">
          <DialogHeader>
            <DialogTitle>Edit Message</DialogTitle>
            <DialogDescription>{props.description}</DialogDescription>
          </DialogHeader>
          <MessageEditor
            data={props.value ? JSON.parse(props.value) : {}}
            onChange={(data) => {
              props.onChange(JSON.stringify(data));
            }}
          />
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
};

const MessageFormatWidget = function (props: WidgetProps) {
  return (
    <MessageFormatEditor
      data={props.value ? JSON.parse(props.value) : {}}
      onChange={(data) => {
        props.onChange(JSON.stringify(data));
      }}
    />
  );
};

const DurationWidget = function (props: WidgetProps & { multiplier: number }) {
  return (
    <Input
      type="number"
      value={Number(props.value) / props.multiplier}
      onChange={(e) =>
        props.onChange(Number(e.target.value) * props.multiplier)
      }
    />
  );
};

const DurationSecondsWidget = function (
  props: WidgetProps & { multiplier?: number }
) {
  return (
    <DurationWidget
      {...props}
      multiplier={1000000000 * (props.multiplier ?? 1)}
    />
  );
};

const DurationMinutesWidget = function (props: WidgetProps) {
  return <DurationSecondsWidget multiplier={60} {...props} />;
};

const DurationHoursWidget = function (props: WidgetProps) {
  return <DurationSecondsWidget multiplier={60 * 60} {...props} />;
};

const DurationDaysWidget = function (props: WidgetProps) {
  return <DurationSecondsWidget multiplier={24 * 60 * 60} {...props} />;
};

const CheckboxWidget = function (props: WidgetProps) {
  return (
    <Switch
      checked={!!props.value}
      onCheckedChange={() => props.onChange(!props.value)}
    />
  );
};

function ArrayFieldTemplate(props: ArrayFieldTemplateProps) {
  return (
    <div>
      <div className="mb-5">
        <div className="text-base font-bold">{props.title}</div>
        <div className="text-sm text-muted-foreground">
          {props.schema.description}
        </div>
      </div>
      <div className="flex flex-col gap-3 mb-3">
        {props.items.map((item) => (
          <Card className="p-5 pl-6 block" key={item.key}>
            <div className="float-right flex items-center gap-3">
              {item.buttonsProps.hasMoveUp && (
                <Button
                  variant="outline"
                  size="icon"
                  onClick={item.buttonsProps.onReorderClick(
                    item.index,
                    item.index - 1
                  )}
                >
                  <ArrowUpIcon />
                </Button>
              )}
              {item.buttonsProps.hasMoveDown && (
                <Button
                  variant="outline"
                  size="icon"
                  onClick={item.buttonsProps.onReorderClick(
                    item.index,
                    item.index + 1
                  )}
                >
                  <ArrowDownIcon />
                </Button>
              )}
              {item.buttonsProps.hasRemove && (
                <Button
                  variant="outline"
                  size="icon"
                  onClick={item.buttonsProps.onDropIndexClick(item.index)}
                >
                  <TrashIcon />
                </Button>
              )}
            </div>

            {item.children}
          </Card>
        ))}
      </div>
      {props.canAdd && (
        <Button type="button" variant="outline" onClick={props.onAddClick}>
          Add Item
        </Button>
      )}
    </div>
  );
}

function ObjectFieldTemplate(props: ObjectFieldTemplateProps) {
  return (
    <div>
      {props.title && (
        <div className="mb-3">
          <div className="text-base font-bold mb-1">{props.title}</div>
          <div className="text-muted-foreground">{props.description}</div>
        </div>
      )}

      <div className="flex flex-col gap-4">
        {props.properties.map((element) => (
          <div key={element.name}>{element.content}</div>
        ))}
      </div>
    </div>
  );
}

function FieldTemplate(props: FieldTemplateProps) {
  const { id, classNames, style, label, help, description, errors, children } =
    props;

  return (
    <div className={classNames} style={style}>
      {!schemaTypeContains(props.schema, "object", "array") && (
        <div className="mb-2">
          <label htmlFor={id} className="block font-bold text-base mb-0.5">
            {label}
          </label>
          <div className="text-muted-foreground text-sm">{description}</div>
        </div>
      )}
      {children}
      {errors}
      {help}
    </div>
  );
}

function LayoutHeaderField(props: FieldProps) {
  return <div className="text-base font-bold bg-red-500">{props.title}</div>;
}

function schemaTypeContains(schema: RJSFSchema, ...types: string[]) {
  for (const type of types) {
    if (schema.type === type || schema.type?.includes(type as any)) {
      return true;
    }
  }
  return false;
}

function FieldErrorTemplate(props: FieldErrorProps) {
  const { errors } = props;

  if (!errors || errors.length === 0) {
    return null;
  }

  return (
    <div className="flex flex-col gap-1 mt-1">
      {errors.map((error, i: number) => {
        return (
          <div key={i} className="text-destructive text-sm font-light">
            {error}
          </div>
        );
      })}
    </div>
  );
}

function SubmitButtonTemplate() {
  return null;
}

const widgets: RegistryWidgetsType = {
  select: SelectWidget,
  channel_select: ChannelSelectWidget,
  role_select: RoleSelectWidget,
  message: MessageWidget,
  message_format: MessageFormatWidget,
  duration_seconds: DurationSecondsWidget,
  duration_minutes: DurationMinutesWidget,
  duration_hours: DurationHoursWidget,
  duration_days: DurationDaysWidget,
  CheckboxWidget: CheckboxWidget,
};

const templates = {
  ObjectFieldTemplate: ObjectFieldTemplate,
  ArrayFieldTemplate: ArrayFieldTemplate,
  FieldErrorTemplate: FieldErrorTemplate,
  LayoutHeaderField: LayoutHeaderField,
  ButtonTemplates: { SubmitButton: SubmitButtonTemplate },
  FieldTemplate: FieldTemplate,
};

export default function SchemaEditor({
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
    <Form
      formData={formData}
      schema={schema}
      uiSchema={uiSchema}
      widgets={widgets}
      templates={templates}
      validator={validator}
      // showErrorList={false}
      onChange={(s) => onChange(s.formData, s.errors)}
      liveValidate
    />
  );
}
