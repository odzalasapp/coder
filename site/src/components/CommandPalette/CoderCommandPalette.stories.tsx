import { Story } from "@storybook/react"
import {
  CoderCommandPalette,
  CoderCommandPaletteProps,
} from "./CoderCommandPalette"

export default {
  title: "components/CoderCommandPalette",
  component: CoderCommandPalette,
}

const Template: Story<CoderCommandPaletteProps> = (
  args: CoderCommandPaletteProps,
) => <CoderCommandPalette {...args} />

export const Example = Template.bind({})
Example.args = {
  isOpen: true,
}
