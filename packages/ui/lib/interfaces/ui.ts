import { ButtonVariants } from "../types/ui";

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  includeDefaultKeywords?: boolean;
  children?: React.ReactNode;
}

export interface IconProps {
  width?: string;
  height?: string;
  fill?: string;
  stroke?: string;
  strokeWidth?: string;
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
}

export interface ErrorType {
  statusCode?: number;
  showRedirectText?: boolean;
  heading?: string;
  statusText?: string;
  message?: string;
  redirectText?: string;
  customPageTitle?: string;
  isFullWidth?: boolean;
  style?: React.CSSProperties;
  hasGetInitialPropsRun?: boolean;
  err?: Error;
  showMessage?: boolean;
  showStatusText?: boolean;
  isModal?: boolean;
  redirectLink?: string;
  children?: React.ReactNode;
  [key: string]: any;
}

export interface ConnectWalletButtonProps {
  className?: string;
  text?: string;
  variant?: ButtonVariants;
  onConnect?: () => void;
  renderContent?: (props: {
    walletConnected: boolean;
    isWrongNetwork: boolean;
    address: string | null;
    text: string;
  }) => React.ReactNode;
}
