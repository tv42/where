(require 'go-mode)

(defcustom go-where-command "go-where"
  "The 'go-where' command."
  :type 'string
  :group 'go)

(define-error 'go-where-error "go-where failed")

(defun go-where--call (lookup)
  "Call go-where for lookup."
  (with-temp-buffer
    (let ((coding-system-for-read 'utf-8)
	  (coding-system-for-write 'utf-8))
      (let ((status (call-process go-where-command nil t nil
				  "--" lookup)))
	(let ((output (buffer-substring-no-properties (point-min) (point-max))))
	  (if (not (= status 0))
	      (signal 'go-where-error output)
	    (replace-regexp-in-string "\n\\'" "" output)))))))

;;;###autoload
(defun go-where (lookup &optional other-window)
  "Find Go identifier described in LOOKUP and find that location."
  (interactive "MLookup (importpath#ident): ")

  ;; TODO: maybe support setting GOOS, GOARCH, build tags?
  (condition-case err
      (let ((file (go-where--call lookup)))
	(push-mark)
	(ring-insert find-tag-marker-ring (point-marker))
	(godef--find-file-line-column file other-window))
    (file-error (message "Could not run go-where binary"))
    (go-where-error (message (cdr err)))))

(provide 'go-where)
