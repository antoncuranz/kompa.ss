import { create, StateCreator } from "zustand"

type PrivacySlice = {
  privacyMode: boolean
  togglePrivacyMode: () => void
}

const createPrivacySlice: StateCreator<
  PrivacySlice,
  [],
  [],
  PrivacySlice
> = set => ({
  privacyMode: false,
  togglePrivacyMode: () => set(state => ({ privacyMode: !state.privacyMode })),
})

export const useStore = create<PrivacySlice>()((...a) => ({
  ...createPrivacySlice(...a),
}))
