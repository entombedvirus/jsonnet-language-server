;;; lsp -- Summary
;; Development lsp registration for Emacs lsp-mode.
;;; Commentary:
;;; Code:
(require 'lsp-mode)

(defcustom lsp-jsonnet-executable "jsonnet-language-server"
  (add-to-list 'lsp-language-id-configuration '(jsonnet-mode . "jsonnet"))
  "Command to start the Jsonnet language server."
  :group 'lsp-jsonnet
  :risky t
  :type 'file)
(lsp-register-client
 (make-lsp-client
  :new-connection (lsp-stdio-connection (lambda () lsp-jsonnet-executable))
  :activation-fn (lsp-activate-on "jsonnet")
  :server-id 'jsonnet))

(provide 'lsp)
;;; lsp.el ends here
