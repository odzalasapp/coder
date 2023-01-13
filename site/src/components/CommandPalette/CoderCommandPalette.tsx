import { FC, useState } from "react"
import { useSearchParams } from "react-router-dom"
import {
  CommandPalette,
  useCommandPalette,
  Option,
  Filter,
} from "./CommandPalette"

export interface CoderCommandPaletteProps {
  isOpen: boolean
}

export const CoderCommandPalette: FC<CoderCommandPaletteProps> = ({
  isOpen,
}) => {
  return (
    <CommandPalette
      isOpen={isOpen}
      defaultView="listRepos"
      views={{
        listRepos: <ListReposView />,
        listWorkspaces: <ListWorkspacesView />,
      }}
    />
  )
}

const ListReposView = () => {
  const { setActiveOption, setView, state } = useCommandPalette()
  const [filter, setFilter] = useState("")
  const [searchParams, setSearchParams] = useSearchParams()

  return (
    <div>
      <Filter
        value={filter}
        onChange={setFilter}
        placeholder="Select a repository to open..."
      />
      <div>
        <Option
          isActive={state.activeOptionIndex === 0}
          onActive={() => setActiveOption(0)}
          onSelect={() => {
            setView("listWorkspaces")
            setSearchParams({ ...searchParams, repo: "coder/docs" })
          }}
        >
          coder/docs
        </Option>
        <Option
          isActive={state.activeOptionIndex === 1}
          onActive={() => setActiveOption(1)}
          onSelect={() => {
            setView("listWorkspaces")
            setSearchParams({ ...searchParams, repo: "coder/coder" })
          }}
        >
          coder/coder
        </Option>
        <Option
          isActive={state.activeOptionIndex === 2}
          onActive={() => setActiveOption(2)}
          onSelect={() => {
            setView("listWorkspaces")
            setSearchParams({ ...searchParams, repo: "coder/coder.com" })
          }}
        >
          coder/coder.com
        </Option>
      </div>
    </div>
  )
}

const ListWorkspacesView = () => {
  const { setActiveOption, state } = useCommandPalette()
  const [filter, setFilter] = useState("")
  const [searchParams] = useSearchParams()

  return (
    <div>
      <Filter
        value={filter}
        onChange={setFilter}
        placeholder="Select a workspace..."
      />
      <div>
        <Option
          isActive={state.activeOptionIndex === 0}
          onActive={() => setActiveOption(0)}
          onSelect={() => {
            alert(`Workspace: bruno-dev, repo: ${searchParams.get("repo")}`)
          }}
        >
          Bruno-dev
        </Option>
        <Option
          isActive={state.activeOptionIndex === 1}
          onActive={() => setActiveOption(1)}
          onSelect={() => {
            alert(`Workspace: bruno-brazil, repo: ${searchParams.get("repo")}`)
          }}
        >
          Bruno Brazil
        </Option>
      </div>
    </div>
  )
}
