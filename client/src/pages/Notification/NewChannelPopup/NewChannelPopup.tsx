export interface NewChannelPopupProps {
  isDisplayed: boolean;
  setIsDisplayed: (arg: boolean) => void;
  onSubmit: () => void;
}

export interface NewChannelPopup {
  popup: ({
    isDisplayed,
    setIsDisplayed,
  }: NewChannelPopupProps) => React.ReactNode;
}
