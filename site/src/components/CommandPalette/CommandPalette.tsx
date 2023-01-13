import Dialog from "@material-ui/core/Dialog"
import { makeStyles } from "@material-ui/core/styles"
import { Stack } from "components/Stack/Stack"
import {
  createContext,
  FC,
  PropsWithChildren,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react"
import SearchIcon from "@material-ui/icons/SearchOutlined"

export interface CommandPaletteProps {
  isOpen: boolean
  defaultView: string
  views: Record<string, JSX.Element>
}

interface CommandPaletteState {
  status: "idle" | "loading"
  activeOptionIndex: number
  view: string
}

interface CommandPaletteContextValue {
  state: CommandPaletteState
  setActiveOption: (index: number) => void
  setView: (view: string) => void
}

const CommandPaletteContext = createContext<
  CommandPaletteContextValue | undefined
>(undefined)

export const CommandPalette: FC<CommandPaletteProps> = ({
  isOpen,
  defaultView,
  views,
}) => {
  const styles = useStyles()
  const [state, setState] = useState<CommandPaletteState>({
    status: "idle",
    view: defaultView,
    activeOptionIndex: 0,
  })

  useEffect(() => {
    const keydownHandler = (event: KeyboardEvent) => {
      switch (event.key) {
        case "ArrowDown":
          setState((state) => {
            return assignState(state, {
              activeOptionIndex: state.activeOptionIndex + 1,
            })
          })
          return

        case "ArrowUp":
          setState((state) => {
            if (state.activeOptionIndex === 0) {
              return assignState(state, { activeOptionIndex: 0 })
            }

            return assignState(state, {
              activeOptionIndex: state.activeOptionIndex - 1,
            })
          })
          return
      }
    }

    document.addEventListener("keydown", keydownHandler)

    return () => {
      document.removeEventListener("keydown", keydownHandler)
    }
  }, [])

  const setActiveOption = (index: number) => {
    setState((state) => assignState(state, { activeOptionIndex: index }))
  }

  const setView = (view: string) => {
    setState((state) =>
      assignState(state, { view, activeOptionIndex: undefined }),
    )
  }

  return (
    <Dialog open={isOpen} PaperProps={{ className: styles.dialogPaper }}>
      <CommandPaletteContext.Provider
        value={{ state, setActiveOption, setView }}
      >
        {views[state.view]}
      </CommandPaletteContext.Provider>
    </Dialog>
  )
}

export const useCommandPalette = (): CommandPaletteContextValue => {
  const context = useContext(CommandPaletteContext)

  if (!context) {
    throw new Error(
      "useCommandPalette only should be used inside of <CommandPalette />",
    )
  }

  return context
}

type OptionProps = PropsWithChildren<{
  onActive: () => void
  onSelect: () => void
  isActive: boolean
}>

export const Option: FC<OptionProps> = ({
  children,
  isActive,
  onSelect,
  onActive,
}) => {
  const styles = useStyles()
  const optionRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (isActive && optionRef.current) {
      optionRef.current.focus()
    }
  }, [isActive])

  return (
    <div
      ref={optionRef}
      className={styles.option}
      role="button"
      tabIndex={-1}
      onMouseEnter={onActive}
      onFocus={onActive}
      onClick={onSelect}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          onSelect()
        }
      }}
    >
      {children}
    </div>
  )
}

interface FilterInputProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

export const Filter: FC<FilterInputProps> = ({
  value,
  onChange,
  placeholder,
}) => {
  const styles = useStyles()
  return (
    <Stack
      direction="row"
      spacing={0}
      alignItems="center"
      className={styles.filter}
    >
      <div className={styles.filterIcon}>
        <SearchIcon />
      </div>
      <input
        className={styles.filterInput}
        type="text"
        value={value}
        onChange={(e) => onChange(e.currentTarget.value)}
        placeholder={placeholder}
      />
    </Stack>
  )
}

const assignState = (
  state: CommandPaletteState,
  partialState: Partial<CommandPaletteState>,
) => {
  return {
    ...state,
    ...partialState,
  }
}

const useStyles = makeStyles((theme) => ({
  dialogPaper: {
    width: "100%",
    maxWidth: theme.spacing(87),
  },
  option: {
    padding: theme.spacing(1.5, 2),
    cursor: "pointer",
    outline: "none",

    "&:focus": {
      background: theme.palette.background.paperLight,
    },
  },
  filter: {
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  filterIcon: {
    padding: theme.spacing(0, 2),
    flexShrink: 0,
    lineHeight: 0,

    "& svg": {
      width: 20,
      height: 20,
    },
  },
  filterInput: {
    width: "100%",
    border: 0,
    background: "transparent",
    fontSize: 16,
    padding: theme.spacing(2, 2, 2, 0),
    outline: "none",
  },
}))
