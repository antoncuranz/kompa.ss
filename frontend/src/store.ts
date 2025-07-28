import {create, StateCreator} from 'zustand'
import {Account, Credit} from "@/types.ts";

type SelectionSlice = {
  currentAccount: Account | null;
  creditSelection: Credit | null;
  exchangeSelection: number | null;
  setCurrentAccount: (value: SelectionSlice["currentAccount"]) => void;
  setCreditSelection: (value: SelectionSlice["creditSelection"], accounts: Account[]) => void;
  setExchangeSelection: (value: SelectionSlice["exchangeSelection"]) => void;
};

const createSelectionSlice: StateCreator<
  SelectionSlice,
  [],
  [],
  SelectionSlice
> = (set) => ({
  currentAccount: null,
  creditSelection: null,
  exchangeSelection: null,
  setCurrentAccount: (currentAccount) => set(() => ({ currentAccount: currentAccount })),
  setCreditSelection: (value, accounts) => set((state) => ({
    creditSelection: value, currentAccount: accounts.find(acct => acct.id == value?.account_id) ?? state.currentAccount
  })),
  setExchangeSelection: (exchangeSelection) => set(() => ({ exchangeSelection: exchangeSelection })),
})

type TransactionAmountSlice = {
  changedTransactionAmounts: {[id: number]: number|null};
  clearTransactionAmounts: () => void,
  clearTransactionAmount: (id: number) => void,
  putTransactionAmount: (id: number, amount: number|null) => void,
};

const createTransactionAmountSlice: StateCreator<
  TransactionAmountSlice,
  [],
  [],
  TransactionAmountSlice
> = (set) => ({
  changedTransactionAmounts: {},
  clearTransactionAmounts: () => set(() => ({ changedTransactionAmounts: {} })),
  clearTransactionAmount: (id) => set((state) => {
    const newChangedTransactionAmounts = state.changedTransactionAmounts
    delete newChangedTransactionAmounts[id]
    return ({ changedTransactionAmounts: newChangedTransactionAmounts })
  }),
  putTransactionAmount: (id, amount) => set((state) => ({ changedTransactionAmounts: { ...state.changedTransactionAmounts, [id]: amount } })),
})

type FilterSlice = {
  filterMissingAmountEur: boolean;
  setFilterMissingAmountEur: (enabled: boolean) => void;
};

const createFilterSlice: StateCreator<
  FilterSlice,
  [],
  [],
  FilterSlice
> = (set) => ({
  filterMissingAmountEur: false,
  setFilterMissingAmountEur: (enabled) => set(() => ({ filterMissingAmountEur: enabled })),
})

type PrivacySlice = {
  privacyMode: boolean;
  togglePrivacyMode: () => void,
};

const createPrivacySlice: StateCreator<
  PrivacySlice,
  [],
  [],
  PrivacySlice
> = (set) => ({
  privacyMode: false,
  togglePrivacyMode: () => set(state => ({ privacyMode: !state.privacyMode })),
})

export const useStore = create<SelectionSlice & TransactionAmountSlice & FilterSlice & PrivacySlice>()((...a) => ({
  ...createSelectionSlice(...a),
  ...createTransactionAmountSlice(...a),
  ...createPrivacySlice(...a),
  ...createFilterSlice(...a),
}))